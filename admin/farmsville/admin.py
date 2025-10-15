from django.contrib import admin
from django.utils.html import format_html
from django.conf import settings
from django import forms
from .models import Event, ProductName, Photo, Product, ProductClaimed, BlogPost, ContentBlock


@admin.register(Event)
class EventAdmin(admin.ModelAdmin):
    list_display = ['date', 'product_count', 'total_quantity', 'total_claimed']
    list_filter = ['date']
    search_fields = ['date']
    date_hierarchy = 'date'

    def product_count(self, obj):
        return obj.products.count()
    product_count.short_description = 'Products'

    def total_quantity(self, obj):
        return sum(p.qty for p in obj.products.all())
    total_quantity.short_description = 'Total Qty'

    def total_claimed(self, obj):
        total = 0
        for product in obj.products.all():
            total += sum(claim.qty for claim in product.claims.all())
        return total
    total_claimed.short_description = 'Total Claimed'


@admin.register(ProductName)
class ProductNameAdmin(admin.ModelAdmin):
    list_display = ['name', 'product_usage_count']
    search_fields = ['name']
    ordering = ['name']

    def product_usage_count(self, obj):
        return obj.products.count()
    product_usage_count.short_description = 'Times Used'


class PhotoAdminForm(forms.ModelForm):
    upload_file = forms.FileField(required=False, label="Upload Photo")

    class Meta:
        model = Photo
        fields = '__all__'


@admin.register(Photo)
class PhotoAdmin(admin.ModelAdmin):
    form = PhotoAdminForm
    list_display = ['name', 'filename', 'caption', 'photo_preview', 'usage_count']
    search_fields = ['name', 'filename', 'caption']
    fields = ['name', 'caption', 'upload_file', 'filename', 'photo_preview']
    readonly_fields = ['filename', 'photo_preview']

    def usage_count(self, obj):
        return obj.products.count()
    usage_count.short_description = 'Used By'

    def photo_preview(self, obj):
        if obj.filename:
            full_url = f"{settings.PHOTOS_URL}/dev_photos/{obj.filename}"
            return format_html(
                '<img src="{}" style="max-width: 400px; max-height: 400px; border: 1px solid #ddd; border-radius: 4px;" />',
                full_url
            )
        return "No photo"
    photo_preview.short_description = 'Preview'

    def save_model(self, request, obj, form, change):
        if settings.IS_PRODUCTION:
            # TODO: Implement photo upload for production (e.g., S3, CDN)
            # For now, just save the model without handling photo uploads
            super().save_model(request, obj, form, change)
            return

        upload_file = form.cleaned_data.get('upload_file')

        if upload_file:
            photos_base = settings.BASE_DIR.parent / 'data' / 'photos' / 'dev_photos'

            # Save new photo
            filename = upload_file.name
            filepath = photos_base / filename

            with open(filepath, 'wb+') as destination:
                for chunk in upload_file.chunks():
                    destination.write(chunk)

            obj.filename = filename

        super().save_model(request, obj, form, change)

    def delete_model(self, request, obj):
        print("delete_model was called")
        if settings.IS_PRODUCTION:
            # TODO: Implement photo deletion for production (e.g., S3, CDN)
            # For now, just delete the model without handling photo file deletion
            super().delete_model(request, obj)
            return

        if obj.filename:
            photos_base = settings.BASE_DIR.parent / 'data' / 'photos' / 'dev_photos'
            filepath = photos_base / obj.filename
            print(f"Deleting file: {filepath}")

            # Delete the photo file if it exists
            if filepath.exists():
                filepath.unlink()
                print("File deleted successfully")
            else:
                print("File does not exist")

        super().delete_model(request, obj)

    def delete_queryset(self, request, queryset):
        print("delete_queryset was called")
        if settings.IS_PRODUCTION:
            # TODO: Implement photo deletion for production (e.g., S3, CDN)
            # For now, just delete the queryset without handling photo file deletion
            super().delete_queryset(request, queryset)
            return

        photos_base = settings.BASE_DIR.parent / 'data' / 'photos' / 'dev_photos'

        # Delete all photo files before deleting the records
        for obj in queryset:
            if obj.filename:
                filepath = photos_base / obj.filename
                print(f"Deleting file: {filepath}")
                if filepath.exists():
                    filepath.unlink()
                    print("File deleted successfully")
                else:
                    print("File does not exist")

        super().delete_queryset(request, queryset)


class ProductClaimedInline(admin.TabularInline):
    model = ProductClaimed
    extra = 0
    fields = ['datetime', 'user_name', 'qty', 'notes']
    readonly_fields = []
    can_delete = True


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ['product_name', 'event', 'qty', 'remaining', 'claimed_qty',
                    'availability_status', 'has_photo']
    list_filter = ['event', 'product_name']
    search_fields = ['product_name__name', 'notes']
    autocomplete_fields = ['product_name', 'photo']
    date_hierarchy = 'event__date'
    fields = ['event', 'product_name', 'qty', 'remaining', 'notes', 'photo', 'photo_preview']
    readonly_fields = ['photo_preview']
    inlines = [ProductClaimedInline]

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.select_related('event', 'product_name', 'photo').prefetch_related('claims')

    def claimed_qty(self, obj):
        return sum(claim.qty for claim in obj.claims.all())
    claimed_qty.short_description = 'Claimed'

    def availability_status(self, obj):
        if obj.remaining == 0:
            color = 'red'
            status = 'Out of Stock'
        elif obj.remaining < obj.qty * 0.25:
            color = 'orange'
            status = 'Low Stock'
        else:
            color = 'green'
            status = 'Available'
        return format_html(
            '<span style="color: {}; font-weight: bold;">{}</span>',
            color, status
        )
    availability_status.short_description = 'Status'

    def has_photo(self, obj):
        return bool(obj.photo)
    has_photo.short_description = 'Photo'
    has_photo.boolean = True

    def photo_preview(self, obj):
        if obj.photo and obj.photo.filename:
            full_url = f"{settings.PHOTOS_URL}/dev_photos/{obj.photo.filename}"
            return format_html(
                '<img src="{}" style="max-width: 400px; max-height: 400px; border: 1px solid #ddd; border-radius: 4px;" />',
                full_url
            )
        return "No photo"
    photo_preview.short_description = 'Current Photo'


@admin.register(ProductClaimed)
class ProductClaimedAdmin(admin.ModelAdmin):
    list_display = ['user_name', 'product_name_display', 'event_display',
                    'qty', 'datetime', 'has_notes']
    list_filter = ['product__event', 'datetime', 'user_name']
    search_fields = ['user_name', 'product__product_name__name', 'notes']
    date_hierarchy = 'datetime'
    autocomplete_fields = ['product']
    fields = ['product', 'datetime', 'user_name', 'qty', 'notes']

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.select_related('product', 'product__event', 'product__product_name')

    def product_name_display(self, obj):
        return obj.product.product_name.name
    product_name_display.short_description = 'Product'
    product_name_display.admin_order_field = 'product__product_name__name'

    def event_display(self, obj):
        return obj.product.event.date
    event_display.short_description = 'Event Date'
    event_display.admin_order_field = 'product__event__date'

    def has_notes(self, obj):
        return bool(obj.notes)
    has_notes.short_description = 'Notes'
    has_notes.boolean = True


class ContentBlockInline(admin.StackedInline):
    model = ContentBlock
    extra = 1
    fields = ['block_type', 'order', 'text_content', 'photo', 'youtube_url']
    autocomplete_fields = ['photo']

    class Media:
        css = {
            'all': ('admin/css/forms.css',)
        }


@admin.register(BlogPost)
class BlogPostAdmin(admin.ModelAdmin):
    list_display = ['title', 'is_published', 'created_at', 'block_count']
    list_filter = ['is_published', 'created_at']
    search_fields = ['title']
    date_hierarchy = 'created_at'
    fields = ['title', 'is_published']
    inlines = [ContentBlockInline]

    def block_count(self, obj):
        return obj.content_blocks.count()
    block_count.short_description = 'Content Blocks'

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.prefetch_related('content_blocks')


@admin.register(ContentBlock)
class ContentBlockAdmin(admin.ModelAdmin):
    list_display = ['blog_post', 'block_type', 'order', 'content_preview']
    list_filter = ['block_type', 'blog_post']
    search_fields = ['blog_post__title', 'text_content']
    autocomplete_fields = ['blog_post', 'photo']
    fields = ['blog_post', 'block_type', 'order', 'text_content', 'photo', 'youtube_url']

    def content_preview(self, obj):
        if obj.block_type == 'text' and obj.text_content:
            preview = obj.text_content[:50]
            return f"{preview}..." if len(obj.text_content) > 50 else preview
        elif obj.block_type == 'photo' and obj.photo:
            return f"Photo: {obj.photo.name}"
        elif obj.block_type == 'youtube' and obj.youtube_url:
            return f"YouTube: {obj.youtube_url}"
        return "No content"
    content_preview.short_description = 'Preview'
