import template from './mail-catcher-index.html.twig';
import './mail-catcher-index.scss';

const { Component, Mixin, Utils } = Shopware;
const { Criteria } = Shopware.Data;

Component.register('mail-catcher-index', {
    template,

    inject: ['repositoryFactory'],
    mixins: [
        Mixin.getByName('listing'),
    ],

    metaInfo() {
        return {
            title: this.$createTitle()
        };
    },

    data() {
        return {
            page: 1,
            limit: 25,
            total: 0,
            repository: null,
            items: null,
            isLoading: true,
            filter: {
                salesChannelId: null,
                customerId: null,
                term: null
            }
        }
    },

    computed: {
        columns() {
            return [
                {
                    property: 'createdAt',
                    dataIndex: 'createdAt',
                    label: 'mail-catcher.list.columns.sentDate',
                    primary: true,
                    routerLink: 'mail.catcher.detail'
                },
                {
                    property: 'subject',
                    dataIndex: 'subject',
                    label: 'mail-catcher.list.columns.subject',
                    allowResize: true,
                    routerLink: 'mail.catcher.detail'
                },
                {
                    property: 'receiver',
                    dataIndex: 'receiver',
                    label: 'mail-catcher.list.columns.receiver',
                    allowResize: true
                }
            ]
        },
        mailCatcherRepository() {
            return this.repositoryFactory.create('sw_mail_catcher');
        }
    },

    methods: {
        getList() {
            this.isLoading = true;

            let criteria = new Criteria();

            if (this.filter.term) {
                criteria.setTerm(this.filter.term);
            }

            criteria.addSorting(Criteria.sort('createdAt', 'DESC'))

            return this.mailCatcherRepository.search(criteria, Shopware.Context.api)
                .then((searchResult) => {
                    this.items = searchResult;
                    this.total = searchResult.total;
                    this.isLoading = false;
                });
        },

        resetFilter() {
            this.filter = {
                salesChannelId: null,
                customerId: null,
                term: null
            };

            this.getList();
        }
    },

    watch: {
        filter: {
            deep: true,
            handler: Utils.debounce(function () {
                this.getList();
            }, 400)
        }
    }
});
