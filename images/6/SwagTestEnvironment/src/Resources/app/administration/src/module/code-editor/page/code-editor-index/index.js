import template from './template.twig';

const { Component, Mixin } = Shopware;

Component.register('code-editor-index', {
    inject: ['codeEditorApiService'],
    template,

    mixins: [
        Mixin.getByName('notification'),
    ],

    shortcuts: {
        'SYSTEMKEY+S': 'saveFile',
    },

    data() {
        return {
            selectedLogFile: null,
            files: [],
            fileContent: '',
        }
    },

    async created() {
        this.files = await this.codeEditorApiService.getFiles();
    },

    methods: {
        async onFileSelected() {
            this.fileContent = '';

            this.fileContent = await this.codeEditorApiService.getFile(this.selectedLogFile);
        },

        async saveFile(){
            if (this.selectedLogFile === null) {
                return;
            }

            await this.codeEditorApiService.saveFile(this.selectedLogFile, this.fileContent);

            this.createNotificationSuccess({
                message: 'File saved successfully'
            });
        }
    }
});
