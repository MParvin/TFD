#!/bin/bash

# TFD Docker Setup Script
# This script helps you set up TFD with Docker quickly

set -e

echo "🚀 TFD (Telegram File Downloader) Docker Setup"
echo "============================================="
echo

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first:"
    echo "   https://docs.docker.com/get-docker/"
    exit 1
fi


# Check if config.yaml exists
if [ ! -f "config.yaml" ]; then
    echo "📝 Creating configuration file..."
    cp config.example.yaml config.yaml
    echo "✅ Configuration file created: config.yaml"
    echo
    echo "⚠️  IMPORTANT: Please edit config.yaml with your bot token and settings:"
    echo "   1. Get a bot token from @BotFather on Telegram"
    echo "   2. Get your Chat ID from @userinfobot on Telegram"  
    echo "   3. Edit config.yaml with your values"
    echo
    read -p "Press Enter to continue after editing config.yaml..."
else
    echo "✅ Configuration file already exists: config.yaml"
fi

# Create data directory if it doesn't exist
if [ ! -d "data" ]; then
    echo "📁 Creating data directory..."
    mkdir -p data
    echo "✅ Data directory created: ./data"
fi

# Check if the configuration has been edited
if grep -q "YOUR_BOT_TOKEN_HERE" config.yaml; then
    echo "⚠️  WARNING: It looks like you haven't configured your bot token yet!"
    echo "   Please edit config.yaml and replace 'YOUR_BOT_TOKEN_HERE' with your actual bot token."
    echo
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "👋 Exiting. Please configure your bot token and run this script again."
        exit 1
    fi
fi

echo
echo "🚀 Starting TFD bot..."
docker compose up -d --build

echo
echo "✅ TFD bot is now running!"
echo
echo "📊 Useful commands:"
echo "   View logs:     docker compose logs -f"
echo "   Stop bot:      docker compose down"  
echo "   Restart bot:   docker compose restart"
echo "   Update image:  docker compose pull && docker compose up -d"
echo
echo "📁 Downloaded files will be saved to: ./data/"
echo "📝 Configuration file: ./config.yaml"
echo
echo "🎉 Setup complete! Your bot should now be responding to messages."
echo "   Send a photo, video, or document to your bot to test it!"