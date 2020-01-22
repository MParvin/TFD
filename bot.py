#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Fri Jul  4 12:03:00 2018
@author: mparvin
"""

import sys
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
import logging
# Use configparser to use config file
import configparser
# path and mkdirs to create download folders if not exists
from os import path
from os import makedirs
# Time used for file names
import time

# Validator for check is URL valid
import validators
# Use Path to extract file extension
from pathlib import Path
# Use shutil and requests to download file
import shutil, requests

config = configparser.ConfigParser()
config.read('config')
### Get admin chat_id from config file
### For more security replies only send to admin chat_id
adminCID  = config['SecretConfig']['adminCID']
photoPath = config['SecretConfig']['photoDirectoryPath']
videoPath = config['SecretConfig']['videoDirectoryPath']
musicPath = config['SecretConfig']['musicDirectoryPath']
voicePath = config['SecretConfig']['voicesDirectoryPath']
documentPath = config['SecretConfig']['documentDirectoryPath']
otherPath = config['SecretConfig']['otherDirectoryPath']

if not path.isdir(videoPath):
	makedirs(videoPath)
if not path.isdir(photoPath):
	makedirs(photoPath)
if not path.isdir(musicPath):
	makedirs(musicPath)
if not path.isdir(voicePath):
	makedirs(voicePath)


# Enable logging
logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)


# Define a few command handlers. These usually take the two arguments bot and
# update. Error handlers also receive the raised TelegramError object in error.
def start(bot, update):
    """Send a message when the command /start is issued."""
    update.message.reply_text('Hello, This bot is private!!!')


def help(bot, update):
    """Send a message when the command /help is issued."""
    update.message.reply_text('This bot download any file, video or images that recieved.\n To use this bot, send file or forward the post!')


def isValidUser(chat_id):
	usersChatID = config['SecretConfig']['usersChatID'].split(',')
	if chat_id not in usersChatID:
		return True
	else:
		return False

# Download file from telegram or any other sites
def fileDownloader(bot, update, messageID, storeFolder, fileExtension=None, fileID='others', fileName=None, fileURL=None):
	try:
		# File from telegram
		if fileID != 'others':
			fileName = str(int(time.time()))
			filePath = """{}{}.{}""".format(storeFolder, fileName, fileExtension)
			newFile = bot.get_file(fileID)
			newFile.download(filePath)
			update.message.reply_text("Download finished", reply_to_message_id=messageID)
		# File from other sites
		else:
			filePath = str(storeFolder)+str(fileName)
			print("_________________________________________________")
			print(fileURL, filePath)
			print("_________________________________________________")
			with requests.get(fileURL, stream=True) as fileRequest:
				with open(filePath, 'wb') as fileToWrite:
					shutil.copyfileobj(fileRequest.raw, fileToWrite)
					update.message.reply_text("Download finished", reply_to_message_id=messageID)


	except Exception as ex:
		template = "An exception of type {0} occured, Arguments:\n{1!r}\n in line {2}"
		errorMessage = template.format(type(ex).__name__,ex.args,sys.exc_info()[-1].tb_lineno)
		### Send error to admin
		bot.sendMessage(text=errorMessage, chat_id=adminCID)

def checkExtension(bot,fileURL):
	if validators.url(str(fileURL)):
		fileURL = Path(fileURL).name
		fileName = fileURL.split('?')[0]
		fileExtension = fileName.split('.')[-1]
		return fileName, fileExtension
	else:
		return False



def allCheck(bot, update):
	chat_id = update.message.chat_id
	if not isValidUser(chat_id):
		errorMessage = """You cannot use this bot, this is a private bot!!!!\n\n
						To use this bot contact @{}""".format(config['SecretConfig']['supportUser'])
		update.message.reply_text(errorMessage)
	else:
		messageID = update.message.message_id

		# fileName and fileURL will be replace in next steps, we set it to nothing to prevent from errors
		fileName = fileURL = 'nothing'
		# print(update.message)
		if update.message.audio:
			fileID = update.message.audio.file_id
			fileExtension = 'mp3'
			storeFolder = musicPatcheckExtensionh
		elif update.message.voice:
			fileID = update.message.voice.file_id
			fileExtension = 'ogg'
			storeFolder = voicePath
		elif update.message.video:
			fileID = update.message.video.file_id
			fileExtension = 'mp4'
			storeFolder = videoPath
		elif update.message.photo:
			fileID = update.message.photo[-1].file_id
			fileExtension = 'jpg'
			storeFolder = photoPath
		elif update.message.text:
			fileURL = update.message.text
			fileName, fileExtension = checkExtension(bot, fileURL)
			# Use fileID to detect file is for telegram or not
			fileID='others'
		else:
			update.message.reply_text("Only send these files:\n Video, Photo, Audio, Voice Or file URL ;)", reply_to_message_id=messageID)
			return

		if fileExtension in 'mp3 wav':
			storeFolder = musicPath
		elif fileExtension in 'mp4 mpg mkv mov wmv flv':
			storeFolder = videoPath
		elif fileExtension in 'jpg jpeg gif png tif ico ps psd svg':
			storeFolder = photoPath
		elif fileExtension in 'txt doc docx':
			storeFolder = documentPath
		elif fileExtension == 'ogg':
			storeFolder = voicePath
		else:
			storeFolder = otherPath

		# If everything is Ok, a message will be send to user
		update.message.reply_text("Downloading file, Please wait", reply_to_message_id=messageID)

		# Lets download file
		fileDownloader(bot, update, messageID, storeFolder, fileExtension, fileID, fileName, fileURL)




def error(bot, update, error):
    """Log Errors caused by Updates."""
    errorMessage = """Update {} caused error {}""".format(update, error)
    ### Logging errors
    logger.warning()
    ### Send error to admin
    bot.sendMessage(text=errorMessage, chat_id=adminCID)


def main():
    """Start the bot."""
    # Create the EventHandler and pass it your bot's token.
    updater = Updater(config['SecretConfig']['Token'])

    # Get the dispatcher to register handlers
    dp = updater.dispatcher

    # on different commands - answer in Telegram
    dp.add_handler(CommandHandler("start", start))
    dp.add_handler(CommandHandler("help", help))

    # on noncommand i.e message - Download files
    dp.add_handler(MessageHandler([Filters.photo, Filters.video, Filters.audio, Filters.voice, Filters.text], allCheck))

    # log all errors
    dp.add_error_handler(error)

    # Start the Bot
    updater.start_polling()

    # Run the bot until you press Ctrl-C or the process receives SIGINT,
    # SIGTERM or SIGABRT. This should be used most of the time, since
    # start_polling() is non-blocking and will stop the bot gracefully.
    updater.idle()


if __name__ == '__main__':
	main()
