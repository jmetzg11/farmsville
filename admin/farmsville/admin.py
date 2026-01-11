from django.contrib import admin
from django.utils.html import format_html
from django.utils import timezone
from django.conf import settings
from django import forms
from django.shortcuts import redirect
from supabase import create_client
from django.db.models import Count, Sum
from .models import Event, ProductName, Photo, Product, ProductClaimed, BlogPost, ContentBlock, AddProductShortcut


@admin.register(Event)
class EventAdmin(admin.ModelAdmin):
    list_display = ['date', 'product_count', 'total_quantity', 'total_claimed']
    list_filter = ['date']
    search_fields = ['date']
    date_hierarchy = 'date'

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.annotate(
            _product_count=Count('products', distinct=True),
            _total_quantity=Sum('products__qty'),
            _total_claimed=Sum('products__claims__qty')
        )

    def product_count(self, obj):
        return obj._product_count
    product_count.short_description = 'Products'
    product_count.admin_order_field = '_product_count'

    def total_quantity(self, obj):
        return obj._total_quantity or 0
    total_quantity.short_description = 'Total Qty'
    total_quantity.admin_order_field = '_total_quantity'

    def total_claimed(self, obj):
        return obj._total_claimed or 0
    total_claimed.short_description = 'Total Claimed'
    total_claimed.admin_order_field = '_total_claimed'


@admin.register(ProductName)
class ProductNameAdmin(admin.ModelAdmin):
    list_display = ['name', 'product_usage_count']
    search_fields = ['name']
    ordering = ['name']

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.annotate(_product_count=Count('products'))

    def product_usage_count(self, obj):
        return obj._product_count
    product_usage_count.short_description = 'Times Used'
    product_usage_count.admin_order_field = '_product_count'


class PhotoAdminForm(forms.ModelForm):
    upload_file = forms.FileField(required=False, label="Upload Photo")

    class Meta:
        model = Photo
        fields = '__all__'


@admin.register(Photo)
class PhotoAdmin(admin.ModelAdmin):
    form = PhotoAdminForm
    list_display = ['name', 'filename', 'photo_type', 'caption', 'photo_preview', 'usage_count']
    list_filter = ['photo_type']
    search_fields = ['name', 'filename', 'caption']
    fields = ['name', 'photo_type', 'caption', 'upload_file', 'filename', 'photo_preview']
    readonly_fields = ['filename', 'photo_preview']

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.annotate(_usage_count=Count('products'))

    def usage_count(self, obj):
        return obj._usage_count
    usage_count.short_description = 'Used By'
    usage_count.admin_order_field = '_usage_count'

    def photo_preview(self, obj):
        if obj.filename:
            subdir = 'product' if obj.photo_type == Photo.PhotoType.PRODUCT else 'blog'
            full_url = f"{settings.PHOTOS_URL}/{subdir}/{obj.filename}"
            return format_html(
                '<img src="{}" style="max-width: 400px; max-height: 400px; border: 1px solid #ddd; border-radius: 4px;" />',
                full_url
            )
        return "No photo"
    photo_preview.short_description = 'Preview'

    def save_model(self, request, obj, form, change):
        upload_file = form.cleaned_data.get('upload_file')

        if upload_file:
            subdir = 'product' if obj.photo_type == Photo.PhotoType.PRODUCT else 'blog'
            filename = upload_file.name

            if settings.IS_PRODUCTION:
                supabase = create_client(settings.SUPABASE_URL, settings.SUPABASE_SERVICE_ROLE_KEY)
                file_content = upload_file.read()
                bucket_name = 'photos'
                file_path = f"{subdir}/{filename}"

                supabase.storage.from_(bucket_name).upload(
                    file_path,
                    file_content,
                    file_options={"content-type": upload_file.content_type}
                )

                obj.filename = filename
            else:
                photos_base = settings.BASE_DIR.parent / 'data' / 'photos' / subdir
                filepath = photos_base / filename

                with open(filepath, 'wb+') as destination:
                    for chunk in upload_file.chunks():
                        destination.write(chunk)

                obj.filename = filename

        super().save_model(request, obj, form, change)

    def delete_model(self, request, obj):
        if obj.filename:
            subdir = 'product' if obj.photo_type == Photo.PhotoType.PRODUCT else 'blog'

            if settings.IS_PRODUCTION:
                supabase = create_client(settings.SUPABASE_URL, settings.SUPABASE_SERVICE_ROLE_KEY)
                bucket_name = 'photos'
                file_path = f"{subdir}/{obj.filename}"

                supabase.storage.from_(bucket_name).remove([file_path])
            else:
                photos_base = settings.BASE_DIR.parent / 'data' / 'photos' / subdir
                filepath = photos_base / obj.filename

                if filepath.exists():
                    filepath.unlink()

        super().delete_model(request, obj)

    def delete_queryset(self, request, queryset):
        if settings.IS_PRODUCTION:
            supabase = create_client(settings.SUPABASE_URL, settings.SUPABASE_SERVICE_ROLE_KEY)
            bucket_name = 'photos'
            file_paths = []

            for obj in queryset:
                if obj.filename:
                    subdir = 'product' if obj.photo_type == Photo.PhotoType.PRODUCT else 'blog'
                    file_paths.append(f"{subdir}/{obj.filename}")

            if file_paths:
                supabase.storage.from_(bucket_name).remove(file_paths)
        else:
            for obj in queryset:
                if obj.filename:
                    subdir = 'product' if obj.photo_type == Photo.PhotoType.PRODUCT else 'blog'
                    photos_base = settings.BASE_DIR.parent / 'data' / 'photos' / subdir
                    filepath = photos_base / obj.filename
                    if filepath.exists():
                        filepath.unlink()

        super().delete_queryset(request, queryset)


class ProductClaimedInline(admin.TabularInline):
    model = ProductClaimed
    extra = 0
    fields = ['user_name', 'qty', 'notes', 'datetime']
    readonly_fields = ['datetime']
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

    def get_changeform_initial_data(self, request):
        initial = super().get_changeform_initial_data(request)
        if 'event' not in initial:
            latest_event = Event.objects.first()  # Already ordered by -date
            if latest_event:
                initial['event'] = latest_event.pk
        return initial

    def get_form(self, request, obj=None, **kwargs):
        form = super().get_form(request, obj, **kwargs)
        form.base_fields['qty'].help_text = '(Eggs: in dozen)'
        form.base_fields['remaining'].help_text = '(Eggs: in dozen)'
        return form

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
            subdir = 'product' if obj.photo.photo_type == Photo.PhotoType.PRODUCT else 'blog'
            full_url = f"{settings.PHOTOS_URL}/{subdir}/{obj.photo.filename}"
            return format_html(
                '<img src="{}" style="max-width: 400px; max-height: 400px; border: 1px solid #ddd; border-radius: 4px;" />',
                full_url
            )
        return "No photo"
    photo_preview.short_description = 'Current Photo'

    def save_formset(self, request, form, formset, change):
        instances = formset.save(commit=False)
        for instance in instances:
            if isinstance(instance, ProductClaimed) and not instance.pk:
                instance.datetime = timezone.now()
            instance.save()
        formset.save_m2m()


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

    def formfield_for_foreignkey(self, db_field, request, **kwargs):
        if db_field.name == "photo":
            kwargs["queryset"] = Photo.objects.filter(photo_type=Photo.PhotoType.BLOG)
        return super().formfield_for_foreignkey(db_field, request, **kwargs)

    def get_formset(self, request, obj=None, **kwargs):
        formset = super().get_formset(request, obj, **kwargs)
        if obj:
            existing_count = obj.content_blocks.count()
            formset.form.base_fields['order'].initial = existing_count + 1
        else:
            formset.form.base_fields['order'].initial = 1
        return formset

    class Media:
        js = ('admin.js',)


@admin.register(BlogPost)
class BlogPostAdmin(admin.ModelAdmin):
    list_display = ['title', 'is_published', 'created_at', 'block_count']
    list_filter = ['is_published', 'created_at']
    search_fields = ['title']
    date_hierarchy = 'created_at'
    fields = ['title', 'is_published']
    inlines = [ContentBlockInline]

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.annotate(_block_count=Count('content_blocks'))

    def block_count(self, obj):
        return obj._block_count
    block_count.short_description = 'Content Blocks'
    block_count.admin_order_field = '_block_count'


@admin.register(AddProductShortcut)
class AddProductShortcutAdmin(admin.ModelAdmin):
    def changelist_view(self, request, extra_context=None):
        return redirect('admin:farmsville_product_add')

    def has_add_permission(self, request):
        return False

    def has_change_permission(self, request, obj=None):
        return False

    def has_delete_permission(self, request, obj=None):
        return False
