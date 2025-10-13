from django.contrib import admin
from django.utils.html import format_html
from django.conf import settings
from django import forms
from .models import Event, ProductName, Product, ProductClaimed


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


class ProductClaimedInline(admin.TabularInline):
    model = ProductClaimed
    extra = 0
    fields = ['datetime', 'user_name', 'qty', 'notes']
    readonly_fields = []
    can_delete = True


class ProductAdminForm(forms.ModelForm):
    photo_upload = forms.ImageField(required=False, label="Upload Photo")

    class Meta:
        model = Product
        fields = '__all__'
        widgets = {
            'photo_url': forms.HiddenInput(),
        }


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    form = ProductAdminForm
    list_display = ['product_name', 'event', 'qty', 'remaining', 'claimed_qty',
                    'availability_status', 'has_photo']
    list_filter = ['event', 'product_name']
    search_fields = ['product_name__name', 'notes']
    autocomplete_fields = ['product_name']
    date_hierarchy = 'event__date'
    fields = ['event', 'product_name', 'qty', 'remaining', 'notes', 'photo_preview', 'photo_upload']
    readonly_fields = ['photo_preview']
    inlines = [ProductClaimedInline]

    def get_queryset(self, request):
        qs = super().get_queryset(request)
        return qs.select_related('event', 'product_name').prefetch_related('claims')

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
        return bool(obj.photo_url)
    has_photo.short_description = 'Photo'
    has_photo.boolean = True

    def photo_preview(self, obj):
        if obj.photo_url:
            full_url = f"{settings.PHOTOS_URL}/{obj.photo_url}"
            return format_html(
                '<img src="{}" style="max-width: 400px; max-height: 400px; border: 1px solid #ddd; border-radius: 4px;" />',
                full_url
            )
        return "No photo"
    photo_preview.short_description = 'Current Photo'

    def save_model(self, request, obj, form, change):
        if settings.IS_PRODUCTION:
            # TODO: Implement photo upload for production (e.g., S3, CDN)
            # For now, just save the model without handling photo uploads
            super().save_model(request, obj, form, change)
            return

        photo_upload = form.cleaned_data.get('photo_upload')

        if photo_upload:
            photos_base = settings.BASE_DIR.parent.parent / 'data' / 'photos'

            # Delete old photo if exists
            if obj.photo_url:
                old_photo_path = photos_base / obj.photo_url
                if old_photo_path.exists():
                    old_photo_path.unlink()

            # Create directory for this event date
            event_dir = str(obj.event.date)
            event_path = photos_base / event_dir
            event_path.mkdir(parents=True, exist_ok=True)

            # Save new photo
            filename = photo_upload.name
            filepath = event_path / filename

            with open(filepath, 'wb+') as destination:
                for chunk in photo_upload.chunks():
                    destination.write(chunk)

            obj.photo_url = f"{event_dir}/{filename}"

        super().save_model(request, obj, form, change)


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
        """Optimize queries with select_related"""
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
