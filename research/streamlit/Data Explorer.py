import streamlit as st

st.set_page_config(
    page_title="data-explorer",
    page_icon="🌐",
)

def main():
    st.title("🌐 data-explorer")

    # Introduction section
    st.write(
        """
        Welcome to data-explorer – your ultimate solution for uncovering valuable insights from diverse data sources, 
        all through the simplicity of natural language!
        """
    )

    st.subheader("🚀 data-explorer Apps")
    st.markdown(
        """
        1. 📘 [Basic Question and Answering](./Basic_QA)
        """
    )


if __name__ == "__main__":
    main()
