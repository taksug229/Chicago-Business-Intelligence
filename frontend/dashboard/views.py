import os
from django.shortcuts import render
from django.http import HttpResponse
from dotenv import dotenv_values
import psycopg2
import pandas as pd
import plotly.express as px
import plotly.graph_objs as go

dir_path = "dashboard/sql/"
env_vars = dotenv_values(".env")


def Login2Database():
    global env_vars
    dbname = env_vars.get("POSTGRES_DB")
    host = env_vars.get("POSTGRES_HOST")
    user = env_vars.get("POSTGRES_USER")
    password = env_vars.get("POSTGRES_PASSWORD")

    conn = psycopg2.connect(
        dbname=dbname,
        user=user,
        password=password,
        host=host,
        port="5432",
    )
    return conn


def run_query(query: str) -> pd.DataFrame:
    conn = Login2Database()
    cursor = conn.cursor()
    cursor.execute(query)
    records = cursor.fetchall()
    cursor.close()
    conn.close()

    df = pd.DataFrame(records)
    column_names = [desc[0] for desc in cursor.description]
    df.columns = column_names
    return df


def run_query_from_file(file_path: str) -> pd.DataFrame:
    with open(file_path, "r") as file:
        query = file.read()
    df = run_query(query=query)
    return df


def say_hello(request):
    return render(request, "hello.html", {"name": "to Django Microservice"})


def plot_taxi_from_airports(request):
    global dir_path
    query_file = dir_path + "taxis-from-airport.sql"
    df = run_query_from_file(file_path=query_file)
    df["trip_week"] = pd.to_datetime(df["trip_week"])
    a = px.line(
        df,
        x="trip_week",
        y="trip_cnt",
        color="dropoff_zip_code",
        labels={
            "trip_week": "Date",
            "trip_cnt": "Total Trips",
            "dropoff_zip_code": "Destination Zip Codes",
        },
        title="Taxi Trips from Airports Per Week",
    )
    fig = go.Figure(data=a)
    plot_div = fig.to_html(full_html=False)
    return render(request, "plot.html", {"plot_div": plot_div})


def plot_high_ccvi_ca(request):
    global dir_path
    query_file = dir_path + "high-ccvi-ca.sql"
    df = run_query_from_file(file_path=query_file)
    df["trip_week"] = pd.to_datetime(df["trip_week"])
    a = px.line(
        df,
        x="trip_week",
        y="trip_cnt",
        color="community_area_name",
        labels={
            "trip_week": "Date",
            "trip_cnt": "Total Trips (To/From)",
            "community_area_name": "Community Area",
        },
        title="Taxi Trips in High CCVI Community Areas Per Week",
    )
    fig = go.Figure(data=a)
    plot_div = fig.to_html(full_html=False)
    return render(request, "plot.html", {"plot_div": plot_div})


def plot_building_fee(request):
    global dir_path
    query_file = dir_path + "waive-building-fee.sql"
    df = run_query_from_file(file_path=query_file)
    a = px.histogram(
        df,
        x="issue_year",
        y="total_fee",
        color="community_area_name",
        labels={
            "issue_year": "Issue Year",
            "total_fee": "Total Fee",
            "community_area_name": "Community Area",
        },
        title="Building Fee in Community Areas with High Unemployment and Poverty Rate Per Year",
    )
    a.update_xaxes(tickvals=df["issue_year"].unique())
    fig = go.Figure(data=a)
    plot_div = fig.to_html(full_html=False)
    return render(request, "plot.html", {"plot_div": plot_div})


def plot_new_construction_low_income(request):
    global dir_path
    query_file = dir_path + "new-construction-low-income.sql"
    df = run_query_from_file(file_path=query_file)
    a = px.histogram(
        df,
        x="issue_year",
        y="total_fee",
        color="zip_code",
        labels={
            "issue_year": "Issue Year",
            "total_fee": "Total Fee",
            "zip_code": "Zip Code",
        },
        title="New Construction Building Fee in Low Income Zip Codes Per Year",
    )
    a.update_xaxes(tickvals=df["issue_year"].unique())
    fig = go.Figure(data=a)
    plot_div = fig.to_html(full_html=False)
    return render(request, "plot.html", {"plot_div": plot_div})
