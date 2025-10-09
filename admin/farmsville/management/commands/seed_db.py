from django.core.management.base import BaseCommand
from django.utils import timezone
from datetime import timedelta
from farmsville.models import Event, ProductName, Product, ProductClaimed


class Command(BaseCommand):
    help = 'Seeds the database with test data'

    def handle(self, *args, **kwargs):
        # Create an event 5 days from now
        event_date = timezone.now().date() + timedelta(days=5)
        event, created = Event.objects.get_or_create(date=event_date)

        # Product names, photos, and notes
        product_data = [
            ('Eggs', 'chicken.jpg', 'Farm fresh, free-range'),
            ('Potatoes', 'cow_2.jpg', 'Organic, locally grown'),
            ('Tomatoes', 'cow.jpg', 'Vine-ripened heirloom'),
            ('Cookies', 'llama_2.jpeg', 'Freshly baked this morning'),
            ('Carrots', 'llama.jpeg', 'Sweet and crisp'),
        ]

        products = []
        for name, photo, notes in product_data:
            product_name, _ = ProductName.objects.get_or_create(name=name)

            product = Product.objects.create(
                event=event,
                product_name=product_name,
                qty=10,
                remaining=10,
                notes=notes,
                photo_url=f'/dev_photos/{photo}'
            )
            products.append(product)

        # First product: 1 claim (doesn't take all qty)
        ProductClaimed.objects.create(
            product=products[0],
            datetime=timezone.now(),
            user='Jack',
            qty=3,
            notes='Perfect for breakfast!'
        )
        products[0].remaining = 7
        products[0].save()

        # Second product: 2 claims (takes all qty)
        ProductClaimed.objects.create(
            product=products[1],
            datetime=timezone.now(),
            user='Jill',
            qty=6,
            notes='Great for mashing!'
        )
        ProductClaimed.objects.create(
            product=products[1],
            datetime=timezone.now(),
            user='Mary',
            qty=4,
            notes='Making soup tonight'
        )
        products[1].remaining = 0
        products[1].save()
