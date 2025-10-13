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


class Product(models.Model):
    event = models.ForeignKey(Event, on_delete=models.CASCADE, related_name='products')
    product_name = models.ForeignKey(ProductName, on_delete=models.PROTECT, related_name='products')
    qty = models.IntegerField()
    remaining = models.IntegerField()
    notes = models.TextField(blank=True, null=True)
    photo_url = models.URLField(max_length=500, blank=True, null=True)

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
