from django.contrib import admin
from django.urls import path
from . import views

urlpatterns = [
    path("hello/", views.say_hello),
    path(
        "taxis-from-airport/",
        views.plot_taxi_from_airports,
        name="taxi_trips_from_airport",
    ),
    path("high-ccvi-ca/", views.plot_high_ccvi_ca, name="taxi_trips_high_ccvi"),
    path("waive-building-fee/", views.plot_building_fee, name="building_fee"),
    path(
        "new-construction-low-income/",
        views.plot_new_construction_low_income,
        name="new_construction",
    ),
]
