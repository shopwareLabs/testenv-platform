const locale = 'nl-NL';

if (Shopware.Locale.getByName(locale) === false) {
    Shopware.Locale.register(locale, {});
}
