def load_api_kwargs():
    """ Fast API kwargs """
    kwargs = {
            "title": "Data Explorer",
            "description": "A Backend Rest API for Data Explorer",
            "summary": """
                This is an application that provides a backend service for multiple Vision - Language apps.
            """,
            "version": "0.0.1",
            "terms_of_service": "https://www.hifeyinc.com/terms/",
            "contact": {
                "name": "Temiloluwa Adeoti",
                "url": "https://temiloluwa.github.io/",
                "email": "temilolu74@gmail.com",
            },
            "license_info": {
                "name": "Apache 2.0",
                "url": "https://www.apache.org/licenses/LICENSE-2.0.html",
            },
            "openapi_tags": [
                {
                    "name": "Basic QA",
                    "description": "Basic RAG QA",
                },
                {
                    "name": "Advanced QA",
                    "description": "Advance RAG QA",
                    "externalDocs": {
                        "description": "This includes",
                        "url": "https://fastapi.tiangolo.com/",
                    },
                },
            ],
        }
    
    return kwargs
