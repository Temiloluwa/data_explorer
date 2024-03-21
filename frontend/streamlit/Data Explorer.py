import streamlit as st

st.set_page_config(
    page_title="Data Explorer",
    page_icon="🌐",
)

def main():
    st.title("🌐 Data Explorer")

    # Introduction section
    st.write(
        """
        Welcome to Data Explorer – your ultimate solution for uncovering valuable insights from diverse data sources, 
        all through the simplicity of natural language!
        """
    )

    st.subheader("🚀 Data Explorer Apps")
    st.markdown(
        """
        1. 📘 [Basic Question and Answering](./Basic_QA)
        """
    )


if __name__ == "__main__":
    main()
