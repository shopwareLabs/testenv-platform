const { ApiService } = Shopware.Classes;

class CodeEditor extends ApiService {
    constructor(httpClient, loginService, apiEndpoint = '_action/code-editor') {
        super(httpClient, loginService, apiEndpoint);
    }

    getFiles() {
        const apiRoute = `${this.getApiBasePath()}/files`;
        return this.httpClient.get(
            apiRoute,
            {
                headers: this.getBasicHeaders()
            }
        ).then((response) => {
            return ApiService.handleResponse(response);
        });
    }

    getFile(file) {
        const apiRoute = `${this.getApiBasePath()}/file`;
        return this.httpClient.get(
            apiRoute,
            {
                params: {
                    file,
                },
                headers: this.getBasicHeaders()
            }
        ).then((response) => {
            return response.data.content;
        });
    }

    saveFile(file, content) {
        const apiRoute = `${this.getApiBasePath()}/file`;
        return this.httpClient.put(
            apiRoute,
            {
                content
            },
            {
                params: {
                    file
                },
                headers: this.getBasicHeaders()
            }
        ).then((response) => {
            return ApiService.handleResponse(response);
        });
    }
}

export default CodeEditor;
