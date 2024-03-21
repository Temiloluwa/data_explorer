import streamlit as st
from src.shared import Message
from src.basic_qa.main import (post_question, put_ingest_document)
                               

with st.sidebar:
    uploaded_file = st.file_uploader('Upload a PDF Document', type=['pdf'])
    with st.container():
        # Subheader: Instructions
        st.subheader("Instructions")

        # Summary of instructions
        st.write("""
        1. Upload a PDF study material.
        2. The maximum document size is 10 mb.
        3. You are limited to asking 10 questions for this Demo.
        """)

# intialize messages in state
if "messages" not in st.session_state:
    st.session_state.messages = []

def main():
    # display all messages in state
    for message in st.session_state.messages:
        with st.chat_message(message.role):
            st.markdown(message.content)

    # upon user prompt
    if prompt := st.chat_input("What would you like to ask ?"):
        # add user prompt to state
        st.session_state.messages.append(Message(role="user", content=prompt))

        # display chat message
        with st.chat_message("user"):
            st.markdown(prompt)

        # call ask_question endpoint
        with st.chat_message("assistant"):
            stream = post_question(prompt)
            
            # stream response
            response = st.write_stream(stream)

        #add assistant response to state
        st.session_state.messages.append(Message(role="assistant", content=response))


# Run main function when file is uploaded
if uploaded_file is not None:
    put_ingest_document(uploaded_file)
    main()
else:
    st.title("📘Welcome to Basic Question and Answering App")
    st.write("Simply upload your a pdf file 2MB or less, and ask a maximum of 10 questions!")

