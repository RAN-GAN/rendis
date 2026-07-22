from setuptools import setup, find_packages

setup(
    name="rendis",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "websocket-client>=1.6.1",
    ],
    description="Python client for Rendis",
    author="ran-gan",
)
