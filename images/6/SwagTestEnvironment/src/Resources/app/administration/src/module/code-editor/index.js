import './page/code-editor-index';

Shopware.Module.register('code-editor', {
    type: 'plugin',
    name: 'code-editor.title',
    title: 'code-editor.title',
    description: '',
    color: '#303A4F',

    icon: 'default-device-dashboard',

    routes: {
        index: {
            component: 'code-editor-index',
            path: 'index',
        },
    },

    settingsItem: [
        {
            group: 'plugins',
            to: 'code.editor.index',
            icon: 'default-action-settings',
            name: 'code-editor.title'
        }
    ]
});
