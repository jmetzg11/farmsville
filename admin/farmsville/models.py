from django.db import models


class Event(models.Model):
    date = models.DateField()

    def __str__(self):
        return str(self.date)

    class Meta:
        ordering = ['-date']


class ProductName(models.Model):
    name = models.CharField(max_length=200, unique=True)

    def __str__(self):
        return self.name

    class Meta:
        ordering = ['name']


class Photo(models.Model):
    name = models.CharField(max_length=200, unique=True)
    filename = models.CharField(max_length=200, unique=True)
    caption = models.CharField(max_length=500, blank=True, null=True)

    def __str__(self):
        return self.name

    class Meta:
        ordering = ['name']


class Product(models.Model):
    event = models.ForeignKey(Event, on_delete=models.CASCADE, related_name='products')
    product_name = models.ForeignKey(ProductName, on_delete=models.PROTECT, related_name='products')
    qty = models.IntegerField()
    remaining = models.IntegerField()
    notes = models.TextField(blank=True, null=True)
    photo = models.ForeignKey(Photo, on_delete=models.SET_NULL, null=True, blank=True, related_name='products')

    def __str__(self):
        return f"{self.product_name.name} - {self.qty} ({self.event.date})"

    class Meta:
        ordering = ['event', 'product_name']


class ProductClaimed(models.Model):
    product = models.ForeignKey(Product, on_delete=models.CASCADE, related_name='claims')
    datetime = models.DateTimeField()
    user_name = models.CharField(max_length=200)
    qty = models.IntegerField()
    notes = models.TextField(blank=True, null=True)

    def __str__(self):
        return f"{self.user_name} - {self.qty} of {self.product.product_name.name} ({self.datetime})"

    class Meta:
        ordering = ['-datetime']
        verbose_name_plural = "Products Claimed"


class BlogPost(models.Model):
    title = models.CharField(max_length=200)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    is_published = models.BooleanField(default=False)

    def __str__(self):
        return self.title

    class Meta:
        ordering = ['-created_at']


class ContentBlock(models.Model):
    BLOCK_TYPES = [
        ('text', 'Text/Paragraph'),
        ('photo', 'Photo'),
        ('youtube', 'YouTube Video'),
    ]

    blog_post = models.ForeignKey(BlogPost, on_delete=models.CASCADE, related_name='content_blocks')
    block_type = models.CharField(max_length=20, choices=BLOCK_TYPES)
    order = models.IntegerField(default=0)

    # Text content
    text_content = models.TextField(blank=True, null=True, help_text="For text blocks")

    # Photo content
    photo = models.ForeignKey(Photo, on_delete=models.SET_NULL, null=True, blank=True,
                             help_text="For photo blocks")

    # YouTube content
    youtube_url = models.URLField(blank=True, null=True, help_text="Full YouTube URL or video ID")

    def __str__(self):
        return f"{self.blog_post.title} - {self.get_block_type_display()} (Order: {self.order})"

    class Meta:
        ordering = ['blog_post', 'order']
