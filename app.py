from flask import Flask, render_template, request, send_from_directory
import os
import pytesseract as tess
from PIL import Image
from googletrans import Translator
from gtts import gTTS
import textwrap
import PIL.Image
import google.generativeai as genai

app = Flask(__name__)

# Set the Tesseract-OCR path
tess.pytesseract.tesseract_cmd = r'C:/Program Files/Tesseract-OCR/tesseract'

# Configure Google API key
GOOGLE_API_KEY = "----"
genai.configure(api_key=GOOGLE_API_KEY)

# Function to display markdown
def to_markdown(text):
    text = text.replace('•', '  *')
    return textwrap.indent(text, '> ', predicate=lambda _: True)

# Initialize GenerativeModel
model = genai.GenerativeModel('gemini-pro-vision')

# Create a directory to store uploaded images and audio files (if they don't exist)
if not os.path.exists('uploads'):
    os.makedirs('uploads')

@app.route('/')
def index():
    return render_template('/home/index.html')

@app.route('/process', methods=['POST'])
def process():
    if 'user_text' in request.form and request.form['user_text'].strip() != '':
        # Handle text input
        user_text = request.form['user_text']
        selected_language = request.form.get('language')

        # Translate the user-provided text to the selected language
        translator = Translator()
        translated_text = translator.translate(user_text, src='en', dest=selected_language).text

        # Initialize the text-to-speech engine for the selected language
        if selected_language == 'ta':
            tts = gTTS(text=translated_text, lang='ta')
        elif selected_language == 'hi':
            tts = gTTS(text=translated_text, lang='hi')
        else:
            tts = gTTS(text=translated_text, lang='en')

        # Save the audio file
        audio_path = 'uploads/audio_speech.mp3'
        tts.save(audio_path)

        return render_template('/result/index.html', text_tesseract=translated_text, selected_language = selected_language, audio_path_tesseract=audio_path)
    
    elif 'image' in request.files:
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

        # Generate content from image using GenerativeModel
        response = model.generate_content(img)
        generated_text = response.text

        # Determine the selected language
        selected_language = request.form.get('language')

        # Translate the generated text to the selected language
        translator = Translator()
        translated_text_gemini = translator.translate(generated_text, src='en', dest=selected_language).text

        # Translate the extracted text from Tesseract OCR to the selected language
        translated_text_tesseract = translator.translate(extracted_text, src='en', dest=selected_language).text

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

        # Delete the uploaded image after processing (optional)
        os.remove(image_path)

        return render_template('/result/index.html', 
                               text_gemini=translated_text_gemini, 
                               text_tesseract=translated_text_tesseract, 
                               selected_language=selected_language, 
                               audio_path_gemini=audio_path_gemini, 
                               audio_path_tesseract=audio_path_tesseract)

# Serve static files over HTTPS
@app.route('/uploads/<filename>')
def uploaded_file(filename):
    return send_from_directory('uploads', filename)

if __name__ == '__main__':
    # Run the Flask app with HTTPS support
    app.run(ssl_context=('cert.pem', 'key.pem'), host='0.0.0.0', port=5000)
