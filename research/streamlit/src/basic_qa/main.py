from . import Message

def post_question(prompt: str) -> str:
    """Ask a question to the user and return the answer."""
    pass



def put_ingest_document(uploaded_file) -> None:
    progress_text = f"Preparing document {document.name}. Please wait."
    my_bar = st.progress(0, text=progress_text)

    for percent_complete in range(100):
        time.sleep(0.02)
        my_bar.progress(percent_complete + 1, text=progress_text)
    time.sleep(1)
    my_bar.empty()

