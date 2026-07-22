from setuptools import setup, find_packages
import os

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="rendis",
    version="1.0.1",
    packages=find_packages(),
    install_requires=[
        "websocket-client>=1.6.1",
    ],
    description="Python client for Rendis",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="ran-gan",
)
