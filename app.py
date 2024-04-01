from flask import Flask, render_template, request, send_from_directory
import os
import pytesseract as tess
from PIL import Image
from googletrans import Translator
from gtts import gTTS
import textwrap
from PyPDF2 import PdfReader
import PIL.Image
import google.generativeai as genai
from google.generativeai import GenerativeModel, configure
import platform
import os
import pathlib
import textwrap


app = Flask(__name__)

# Configure Port
Port = 5000

# Configure Host
Host = '0.0.0.0'

# Set the Tesseract-OCR path
if platform.system() == 'Windows':
    # Set Tesseract path for Windows
    tess.pytesseract.tesseract_cmd = r'C:/Program Files/Tesseract-OCR/tesseract.exe'
else:
    # Set Tesseract path for Linux
    tess.pytesseract.tesseract_cmd = '/usr/bin/tesseract'

# Configure Google API key
GOOGLE_API_KEY = "AIzaSyBfq78sE-iZ8RsT_kfP93nGi51UyoiVwH8"
genai.configure(api_key=GOOGLE_API_KEY)

# Function to display markdown
def to_markdown(text):
    text = text.replace('•', '  *')
    return textwrap.indent(text, '> ', predicate=lambda _: True)




# Create a directory to store uploaded images and audio files (if they don't exist)
if not os.path.exists('uploads'):
    os.makedirs('uploads')

def extract_text_from_pdf(pdf_file):
    try:
        pdf_reader = PdfReader(pdf_file)
        text = ""
        for page in pdf_reader.pages:
            text += page.extract_text()
        return text
    except Exception as e:
        print("Error extracting text from PDF:", e)
        return None

@app.route('/')
def index():
    return render_template('/home/index.html')

@app.route('/process', methods=['POST'])
def process():
    # print(request.files.keys())
    #Function to Clear Uploads folder
    delete_files_in_uploadsfolder()
    if 'user_text' in request.form and request.form['user_text'].strip() != '':
        print("\n\n********************* Processing Text ***************************\n\n\n")
        # Handle text input
        user_text = request.form['user_text']
        selected_language = request.form.get('language')

        # Printing Selected Language
        print("*** Selected Language ***\n")
        print(selected_language)
        print("\n**** --------------  ****\n\n")

        # Assuming you want to use the 'gemini-pro' model
        model = GenerativeModel('gemini-pro')

        # Generate text for the prompt 
        try:
            # Wrap the generation in a try-except to catch potential errors
            response = model.generate_content(user_text + "Give response in english should not exceeds 100 words.If it is a thing tell about the time, if it is a command tell what the command will do. Tell me about the content in the text")

            # Convert the generated text to markdown format
            generated_text = to_markdown(response.text)
           
        except Exception as e:
            print(f"An error occurred: {e}")
            generated_text = "Error in AI responce"


        # Printing Generated Text
        print("****  Generated Text  ****\n")
        print(generated_text)
        print("\n****  --------------  ****\n\n")

        # Translate the user-provided text to the selected language
        translator = Translator()

        try:
            # Translate the generated text to the selected language
            if generated_text:
                translated_text_gemini = translator.translate(generated_text, src='en', dest=selected_language).text
                if translated_text_gemini is None:
                    translated_text_gemini = "Translation not available"
            else:
                translated_text_gemini = "No text to translate"
        except Exception as e:
            print("Error translating generated text:", e)
            translated_text_gemini = "Error translating generated text"
 


        # Printing generated Translated Text
        print("****  Generated Text translated to ",selected_language," ****\n")
        print(translated_text_gemini)
        print("\n**** -----------------------------------  ****\n\n")


        # Printing Input Text
        print("*** Input Text ***\n")
        print(user_text)
        print("\n*** ---------- ***\n\n")

        try:
            # Translate the text from User
            if user_text:
                translated_text_user = translator.translate(user_text, src='en', dest=selected_language).text
                if translated_text_user is None:
                    translated_text_user = "Translation not available"
            else:
                translated_text_user = "No text to translate"
        except Exception as e:
            print("Error translating extracted text:", e)
            translated_text_user = "Error translating extracted text"
 
        # Printing User Translated Text
        print("****  User Text translated to ",selected_language," ****\n")
        print(translated_text_user)
        print("\n**** -------------------------------  ****\n\n")
        
        # Initialize the text-to-speech engine for the selected language
        if selected_language == 'ta':
            tts_gemini = gTTS(text=translated_text_gemini, lang='ta')
            tts_user = gTTS(text=translated_text_user, lang='ta')
        elif selected_language == 'hi':
            tts_gemini = gTTS(text=translated_text_gemini, lang='hi')
            tts_user = gTTS(text=translated_text_user, lang='hi')
        else:
            tts_gemini = gTTS(text=translated_text_gemini, lang='en')
            tts_user = gTTS(text=translated_text_user, lang='en')

        # Save the audio files
        audio_path_gemini = 'uploads/audio_speech_gemini.mp3'
        audio_path_tesseract = 'uploads/audio_speech_tesseract.mp3'
        tts_gemini.save(audio_path_gemini)
        tts_user.save(audio_path_tesseract)

        # Printing User Audio Path
        print("****  Audio Path to user text in ",selected_language," ****\n")
        print(audio_path_tesseract)
        print("\n**** ---------------------------------  ****\n\n")

        # Printing Generated Audio Path
        print("****  Audio Path to generated text in ",selected_language," ****\n")
        print( audio_path_gemini)
        print("\n**** ---------------------------------  ****\n\n")

        print("***********************************************************************************\n\n\n\n")
        return render_template('/result/index.html', 
                               text_gemini=translated_text_gemini, 
                               text_tesseract=translated_text_user, 
                               selected_language=selected_language, 
                               audio_path_gemini=audio_path_gemini, 
                               audio_path_tesseract=audio_path_tesseract)
    
    if 'image' in request.files:
        print("\n\n********************************* Processing Image ********************************* \n\n\n")

        # Handle image input
        image_file = request.files['image']

        # Save the uploaded image to the 'uploads' directory
        image_path = os.path.join('uploads', image_file.filename)
        image_file.save(image_path)

        # Open the uploaded image for text extraction
        img = Image.open(image_path)
        extracted_text = tess.image_to_string(img)
       
        # Check if extracted text is empty or None
        if not extracted_text:
           extracted_text = "No Text Found In Image"

        print("********* Extracted Text ***********\n")
        print(extracted_text)
        print("\n******** ----------- ***********\n\n")

        # Initialize GenerativeModel
        model = genai.GenerativeModel('gemini-pro-vision')

        # Generate content from image using GenerativeModel
        response = model.generate_content(img)

        # Check if response contains a valid Part
        if hasattr(response, 'text'):
            generated_text = response.text
        else:
            generated_text = "Content generation failed."

       
        print("********** Generated Text **********\n")
        print(generated_text)
        print("\n********* ----------- **********\n\n")
        
        # Determine the selected language
        selected_language = request.form.get('language')
        
        print("******** Seleceted Language ********\n")
        print(selected_language)
        print("\n******* --------------- ********\n\n")

        # Initialize the translator outside conditional blocks
        translator = Translator()

        try:
            # Translate the generated text to the selected language
            if generated_text:
                translated_text_gemini = translator.translate(generated_text, src='en', dest=selected_language).text
                if translated_text_gemini is None:
                    translated_text_gemini = "Translation not available"
            else:
                translated_text_gemini = "No text to translate"
        except Exception as e:
            print("Error translating generated text:", e)
            translated_text_gemini = "Error translating generated text"

        # Printing generated Translated Text
        print("****  Generated Text translated to ",selected_language," ****\n")
        print(translated_text_gemini)
        print("\n**** -----------------------------------  ****\n\n")

        try:
            # Translate the extracted text from Tesseract OCR to the selected language
            if extracted_text:
                translated_text_tesseract = translator.translate(extracted_text, src='en', dest=selected_language).text
                if translated_text_tesseract is None:
                    translated_text_tesseract = "Translation not available"
            else:
                translated_text_tesseract = "No text to translate"
        except Exception as e:
            print("Error translating extracted text:", e)
            translated_text_tesseract = "Error translating extracted text"

        # Printing Tesseract Translated Text
        print("****  Tesseract Text translated to ",selected_language," ****\n")
        print(translated_text_tesseract)
        print("\n**** -------------------------------  ****\n\n")
       
        # Initialize the text-to-speech engine for the selected language
        if selected_language == 'ta':
            tts_gemini = gTTS(text=translated_text_gemini, lang='ta')
            tts_tesseract = gTTS(text=translated_text_tesseract, lang='ta')
        elif selected_language == 'hi':
            tts_gemini = gTTS(text=translated_text_gemini, lang='hi')
            tts_tesseract = gTTS(text=translated_text_tesseract, lang='hi')
        else:
            tts_gemini = gTTS(text=translated_text_gemini, lang='en')
            tts_tesseract = gTTS(text=translated_text_tesseract, lang='en')

        # Save the audio files
        audio_path_gemini = 'uploads/audio_speech_gemini.mp3'
        audio_path_tesseract = 'uploads/audio_speech_tesseract.mp3'
        tts_gemini.save(audio_path_gemini)
        tts_tesseract.save(audio_path_tesseract)

        # Printing User Audio Path
        print("****  Audio Path to user text in ",selected_language," ****\n")
        print(audio_path_tesseract)
        print("\n**** ---------------------------------  ****\n\n")

        # Printing Generated Audio Path
        print("****  Audio Path to generated text in ",selected_language," ****\n")
        print( audio_path_gemini)
        print("\n**** ---------------------------------  ****\n\n")

        # Delete the uploaded image after processing (optional)
        os.remove(image_path)

        print("********************************************************************************\n\n\n\n")
        return render_template('/result/index.html', 
                               text_gemini=translated_text_gemini, 
                               text_tesseract=translated_text_tesseract, 
                               selected_language=selected_language, 
                               audio_path_gemini=audio_path_gemini, 
                               audio_path_tesseract=audio_path_tesseract)
    
    if 'pdf' in request.files:
            print("\n\n********************************** Processing PDF ***************************************\n\n\n")
            pdf_file = request.files['pdf']
            if pdf_file.filename == '':
                return "No PDF file selected"
            
            # Printing PDF Name
            print("****** PDF Filename *******\n")
            print(pdf_file.filename)
            print("\n**** ----------------- ****\n\n")


            # Create the 'uploads' directory if it doesn't exist
            if not os.path.exists('uploads'):
                os.makedirs('uploads')
            pdf_path = os.path.join('uploads', pdf_file.filename)
            pdf_file.save(pdf_path)
            
            # Printing PDF Path
            print("****** PDF Saved Path *******\n")
            print(pdf_path)
            print("\n**** ------------------- ****\n\n")

            
            selected_language = request.form.get('language')

            #Printing Selected Language
            print("****** Selected Language *******\n")
            print(selected_language)
            print("\n**** ---------------------- ****\n\n")

            with open(pdf_path, 'rb') as file:
                extracted_text = extract_text_from_pdf(file)

            if not extracted_text:
               extracted_text = "No Text in PDF"

            #Printing Extracted Text
            print("****** Extracted Text *******\n")
            print(extracted_text)
            print("\n**** ------------------- ****\n\n")

            os.remove(pdf_path)

            user_text = extracted_text

            genai.configure(api_key=GOOGLE_API_KEY)
            # Assuming you want to use the 'gemini-pro' model
            model = GenerativeModel('gemini-pro')

            # Generate text for the prompt 
            try:
                # Wrap the generation in a try-except to catch potential errors
                response = model.generate_content(user_text + "Give response in english should not exceeds 100 words,If it is book tell book title , short narration , author, published year,if it a slogan or quotes tell about it and author of it.")

                # Convert the generated text to markdown format
                generated_text = to_markdown(response.text)
            
            except Exception as e:
                print(f"An error occurred: {e}")
                generated_text = "Error in AI responce"


            # Printing Generated Text
            print("****  Generated Text  ****\n")
            print(generated_text)
            print("\n****  --------------  ****\n\n")

            # Initialize the translator outside conditional blocks
            translator = Translator()

            try:
                # Translate the generated text to the selected language
                if generated_text:
                    translated_text_gemini = translator.translate(generated_text, src='en', dest=selected_language).text
                    if translated_text_gemini is None:
                        translated_text_gemini = "Translation not available"
                else:
                    translated_text_gemini = "No text to translate"
            except Exception as e:
                print("Error translating generated text:", e)
                translated_text_gemini = "Error translating generated text"
    


            # Printing generated Translated Text
            print("****  Generated Text translated to ",selected_language," ****\n")
            print(translated_text_gemini)
            print("\n**** -----------------------------------  ****\n\n")

            try:
                # Translate the text from User
                if user_text:
                    translated_text_user = translator.translate(user_text, src='en', dest=selected_language).text
                    if translated_text_user is None:
                        translated_text_user = "Translation not available"
                else:
                    translated_text_user = "No text to translate"
            except Exception as e:
                print("Error translating extracted text:", e)
                translated_text_user = "Error translating extracted text"
    
            # Printing User Translated Text
            print("****  Extracted Text translated to ",selected_language," ****\n")
            print(translated_text_user)
            print("\n**** -------------------------------  ****\n\n")
            
            # Initialize the text-to-speech engine for the selected language
            if selected_language == 'ta':
                tts_gemini = gTTS(text=translated_text_gemini, lang='ta')
                tts_user = gTTS(text=translated_text_user, lang='ta')
            elif selected_language == 'hi':
                tts_gemini = gTTS(text=translated_text_gemini, lang='hi')
                tts_user = gTTS(text=translated_text_user, lang='hi')
            else:
                tts_gemini = gTTS(text=translated_text_gemini, lang='en')
                tts_user = gTTS(text=translated_text_user, lang='en')

            # Save the audio files
            audio_path_gemini = 'uploads/audio_speech_gemini.mp3'
            audio_path_tesseract = 'uploads/audio_speech_tesseract.mp3'
            tts_gemini.save(audio_path_gemini)
            tts_user.save(audio_path_tesseract)

            # Printing User Audio Path
            print("****  Audio Path to user text in ",selected_language," ****\n")
            print(audio_path_tesseract)
            print("\n**** ---------------------------------  ****\n\n")

            # Printing Generated Audio Path
            print("****  Audio Path to generated text in ",selected_language," ****\n")
            print( audio_path_gemini)
            print("\n**** ---------------------------------  ****\n\n")

            print("*********************************************************************************\n\n\n\n")
            return render_template('/result/index.html', 
                                text_gemini=translated_text_gemini, 
                                text_tesseract=translated_text_user, 
                                selected_language=selected_language, 
                                audio_path_gemini=audio_path_gemini, 
                                audio_path_tesseract=audio_path_tesseract)
        
    else:
        print("In No Data")
        return "No valid input provided"
   
   
# Serve static files over HTTPS
@app.route('/uploads/<filename>')
def uploaded_file(filename):
    uploads_dir = 'uploads'
    return send_from_directory(uploads_dir, filename)

def to_markdown(text):
    text = text.replace('•', '  *')
    return textwrap.indent(text, '> ')

def delete_files_in_uploadsfolder():
    print("\n***** Clearing Uploads Folder *****\n")
    folder_path = 'uploads'
    # Check if the folder path exists
    if not os.path.exists(folder_path):
        print(f"Folder '{folder_path}' does not exist.")
        return
    
    # Iterate over the files in the folder
    for file_name in os.listdir(folder_path):
        file_path = os.path.join(folder_path, file_name)
        # Check if it's a file
        if os.path.isfile(file_path):
            # Delete the file
            os.remove(file_path)
            print(f"Deleted file: {file_path}")

    print("\n***** -------- *****\n\n")


if __name__ == '__main__':
    # Run the Flask app with HTTPS support
    print("\n\n \********************************************* Server Started ***********************************************/ \n")
    print("          |************* Port :",Port," *************** || ************** Host :",Host," **************| \n")
    app.run(ssl_context=('cert.pem', 'key.pem'), host=Host, port=Port)
  