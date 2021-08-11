const { Application } = Shopware;
import CodeEditor from "./api/code-editor";

Application.addServiceProvider('codeEditorApiService', (container) => {
    const initContainer = Application.getContainer('init');

    return new CodeEditor(initContainer.httpClient, container.loginService);
});

import './module/code-editor';
