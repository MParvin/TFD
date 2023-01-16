# TFD
Telegram File Downloader
  This bot download media files(Video, Photo, Audio and Voice) to specific folder.

## Run with Docker
Add Telegram Token and Chat IDs to config file, then:

```bash
docker-compose up -d
```

## Run without Docker:

  1 - Setup Virtualenv
  
     virtualenv .env -p python3
     
  ** After running command if you get this error :
  
     The path python3 (from --python=python3) does not exist
     
   You have not python3 or python3 is not in your path **

    
  2 - Activate Virtualenv
  
    source .env/bin/activate
    
  3 - Install requiremenets
  
    pip install -r requirements.txt
    
  4 - Edit Config file with your file editor(My editor is vim)
    
    vim config
    
   If You want to run this bot on windows, Your folders path in config file will be like this:
   
    c:\Users\MParvin\...
  
  5 - run bot
  
    python bot.py  
    
   You can use
   
      chmod +x bot.py
   then use 
   
    ./bot.py
  
After running bot send your Video, Audio, Voice and Image to this bot for downloading file.


Todo:
  Add comment to bot.py
