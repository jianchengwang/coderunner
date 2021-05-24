#!/bin/bash
docker pull golang:1.15-alpine
docker pull python:3.9.1-alpine
docker pull openjdk:8u232-jdk
docker pull node:lts-alpine
docker pull gcc:latest
docker-compose up -d