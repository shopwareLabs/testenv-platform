import './component/mail-catcher-index';
import './component/mail-catcher-detail';

const {Module} = Shopware;

Module.register('mail-catcher', {
    type: 'plugin',
    title: 'mail-catcher.general.title',
    description: 'mail-catcher.general.description',
    color: '#F88962',
    icon: 'default-device-dashboard',

    routes: {
        index: {
            component: 'mail-catcher-index',
            path: 'index'
        },
        detail: {
            component: 'mail-catcher-detail',
            path: 'detail/:id',
            meta: {
                parentPath: 'mail.catcher.index'
            }
        }
    },

    settingsItem: [{
        group: 'plugins',
        to: 'mail.catcher.index',
        icon: 'default-device-dashboard',
        name: 'mail-catcher.general.title'
    }],
})
