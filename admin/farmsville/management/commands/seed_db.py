from django.core.management.base import BaseCommand
from django.utils import timezone
from datetime import timedelta
from farmsville.models import Event, ProductName, Photo, Product, ProductClaimed


class Command(BaseCommand):
    help = 'Seeds the database with test data'

    def handle(self, *args, **kwargs):
        # Create an event 5 days from now
        event_date = timezone.now().date() + timedelta(days=5)
        event, created = Event.objects.get_or_create(date=event_date)

        # Create photos based on actual files in dev_photos/
        # Format: (name, filename, caption)
        photo_data = [
            ('Chicken Photo', 'chicken.jpg', 'Fresh farm chicken'),
            ('Cow Photo', 'cow.jpg', 'Grass-fed cow'),
            ('Llama Photo', 'llama.jpeg', 'Friendly llama'),
            ('Llama 2 Photo', 'llama_2.jpeg', 'Another llama'),
            ('Rooster Photo', 'rooster.png', 'Morning rooster'),
            ('Apple Photo', 'apple.png', 'Fresh apple'),
        ]

        photos = {}
        for name, filename, caption in photo_data:
            photo, _ = Photo.objects.get_or_create(
                filename=filename,
                defaults={'name': name, 'caption': caption}
            )
            photos[filename] = photo

        # Product names, photos, quantities, and notes
        product_data = [
            ('Chicken', 'chicken.jpg', 12, 'Farm fresh, free-range'),
            ('Cow', 'cow.jpg', 10, 'Grass-fed, organic'),
            ('Llama', 'llama.jpeg', 8, 'Soft and warm'),
            ('Llama 2', 'llama_2.jpeg', 15, 'Another friendly llama'),
            ('Rooster', 'rooster.png', 10, 'Free-range rooster'),
            ('Apple', 'apple.png', 20, 'Crisp and sweet'),
        ]

        products = []
        for name, photo_filename, qty, notes in product_data:
            product_name, _ = ProductName.objects.get_or_create(name=name)

            product = Product.objects.create(
                event=event,
                product_name=product_name,
                qty=qty,
                remaining=qty,
                notes=notes,
                photo=photos[photo_filename]
            )
            products.append(product)

        # First product (Chicken): Jack claims some - partial stock
        ProductClaimed.objects.create(
            product=products[0],
            datetime=timezone.now(),
            user_name='Jack',
            qty=3,
            notes='Perfect for dinner!'
        )
        products[0].remaining = 9
        products[0].save()

        # Second product (Cow): Multiple claims - OUT OF STOCK
        ProductClaimed.objects.create(
            product=products[1],
            datetime=timezone.now(),
            user_name='Jill',
            qty=6,
            notes='Need for the ranch'
        )
        ProductClaimed.objects.create(
            product=products[1],
            datetime=timezone.now(),
            user_name='Mary',
            qty=4
        )
        products[1].remaining = 0
        products[1].save()

        # Third product (Llama): Fully stocked - NO CLAIMS
        # No claims, stays at full stock

        # Fourth product (Llama 2): Sarah claims with note
        ProductClaimed.objects.create(
            product=products[3],
            datetime=timezone.now(),
            user_name='Sarah',
            qty=5,
            notes='So cute and friendly!'
        )
        products[3].remaining = 10
        products[3].save()

        # Fifth product (Rooster): Jill claims without note
        ProductClaimed.objects.create(
            product=products[4],
            datetime=timezone.now(),
            user_name='Jill',
            qty=4
        )
        products[4].remaining = 6
        products[4].save()

        # Sixth product (Apple): Jack claims without note
        ProductClaimed.objects.create(
            product=products[5],
            datetime=timezone.now(),
            user_name='Jack',
            qty=8
        )
        products[5].remaining = 12
        products[5].save()
