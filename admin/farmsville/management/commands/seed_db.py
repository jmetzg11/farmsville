from django.core.management.base import BaseCommand
from django.utils import timezone
from datetime import timedelta
from farmsville.models import Event, ProductName, Photo, Product, ProductClaimed, BlogPost, ContentBlock


class Command(BaseCommand):
    help = 'Seeds the database with test data'

    def handle(self, *args, **kwargs):
        # Create an event 5 days from now
        event_date = timezone.now().date() + timedelta(days=5)
        event = Event.objects.create(date=event_date)

        # Create photos based on actual files in product/
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
            photo = Photo.objects.create(
                name=name,
                filename=filename,
                caption=caption,
                photo_type=Photo.PhotoType.PRODUCT
            )
            photos[filename] = photo

        # Product names, photos, quantities, and notes
        product_data = [
            ('Eggs', 'chicken.jpg', 12, 'Farm fresh, free-range'),
            ('Cow', 'cow.jpg', 10, 'Grass-fed, organic'),
            ('Llama', 'llama.jpeg', 8, 'Soft and warm'),
            ('Llama 2', 'llama_2.jpeg', 15, 'Another friendly llama'),
            ('Rooster', 'rooster.png', 10, 'Free-range rooster'),
            ('Apple', 'apple.png', 20, 'Crisp and sweet'),
        ]

        products = []
        for name, photo_filename, qty, notes in product_data:
            product_name = ProductName.objects.create(name=name)

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

        # Create blog photos
        blog_photo_1 = Photo.objects.create(
            name='Donkey face',
            filename='download.jpeg',
            caption='Smartest donkey on the farm',
            photo_type=Photo.PhotoType.BLOG
        )

        blog_photo_2 = Photo.objects.create(
            name='Cows Blog Photo',
            filename='cows.jpg',
            caption='Cows in the field',
            photo_type=Photo.PhotoType.BLOG
        )

        # First blog post
        blog1 = BlogPost.objects.create(
            title='First Blog Post - Testing Formatting',
            is_published=True
        )

        ContentBlock.objects.create(
            blog_post=blog1,
            block_type=ContentBlock.BlockType.TEXT,
            order=1,
            text_content='''This is the first paragraph with\t\ttabs and    multiple    spaces to test    formatting.
Let\'s see how this renders on the frontend!

This paragraph has multiple lines
    and some indentation
        with even more indentation here
And back to normal.

Here\'s another section with\ttabs\tbetween\twords and some more text to make it longer so we can see how it wraps and displays across different screen sizes.'''
        )

        ContentBlock.objects.create(
            blog_post=blog1,
            block_type=ContentBlock.BlockType.PHOTO,
            order=2,
            photo=blog_photo_1
        )

        ContentBlock.objects.create(
            blog_post=blog1,
            block_type=ContentBlock.BlockType.TEXT,
            order=3,
            text_content='''This is another paragraph after the photo. More content here to test the layout and see how everything flows together.

Multiple paragraphs within a single content block!
    With some indentation for testing.

And even more text to see how this all renders. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.'''
        )

        ContentBlock.objects.create(
            blog_post=blog1,
            block_type=ContentBlock.BlockType.YOUTUBE,
            order=4,
            youtube_url='rmojioqOlzk'
        )

        # Second blog post
        blog2 = BlogPost.objects.create(
            title='Second Blog Post - More Content',
            is_published=True
        )

        ContentBlock.objects.create(
            blog_post=blog2,
            block_type=ContentBlock.BlockType.TEXT,
            order=1,
            text_content='''First text block of the second blog post with lots of content to test.

This has newlines
    and indentation
        and nested indentation
Back to the start again.

More fluff text here\twith\ttabs\tsprinkled\tthroughout the content to see how they render.'''
        )

        ContentBlock.objects.create(
            blog_post=blog2,
            block_type=ContentBlock.BlockType.TEXT,
            order=2,
            text_content='''Second text block with more information and testing.

Paragraph one of this block.

Paragraph two with    extra    spaces    between    words.

Paragraph three with normal text flow to see the difference.'''
        )

        ContentBlock.objects.create(
            blog_post=blog2,
            block_type=ContentBlock.BlockType.PHOTO,
            order=3,
            photo=blog_photo_2
        )

        ContentBlock.objects.create(
            blog_post=blog2,
            block_type=ContentBlock.BlockType.YOUTUBE,
            order=4,
            youtube_url='C3l27fFeSPk'
        )

        ContentBlock.objects.create(
            blog_post=blog2,
            block_type=ContentBlock.BlockType.TEXT,
            order=5,
            text_content='''Final paragraph wrapping up the second blog post.

This includes multiple paragraphs as well!

    Some indented text here
    And more indented text

Final thoughts with\ttabs\tand    spaces    mixed together to really test the rendering capabilities of the frontend.'''
        )
