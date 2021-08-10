<?php declare(strict_types=1);

namespace SwagDemoProducts;

use Shopware\Core\Content\Product\Aggregate\ProductVisibility\ProductVisibilityDefinition;
use Shopware\Core\Defaults;
use Shopware\Core\Framework\Context;
use Shopware\Core\Framework\DataAbstractionLayer\Search\Criteria;
use Shopware\Core\Framework\DataAbstractionLayer\Search\Filter\EqualsFilter;
use Shopware\Core\Framework\DataAbstractionLayer\Search\Query\ScoreQuery;
use Shopware\Core\Framework\Plugin;
use Shopware\Core\Framework\Plugin\Context\InstallContext;

class SwagDemoProducts extends Plugin
{
    public function install(InstallContext $installContext): void
    {
        $c = Context::createDefaultContext();

        $salesChannelId = $this->container->get('sales_channel.repository')
            ->searchIds((new Criteria())->addFilter(new EqualsFilter('typeId', Defaults::SALES_CHANNEL_TYPE_STOREFRONT)), $c)
            ->firstId();

        $categoryCriteria = new Criteria();
        $categoryCriteria->addQuery(new ScoreQuery(new EqualsFilter('name', 'Industrial & Books'), 1000));
        $categoryCriteria->setLimit(2);

        $categoryId = $this->container->get('category.repository')
            ->searchIds($categoryCriteria, $c)
            ->firstId();

        $taxId = $this->container->get('tax.repository')
            ->searchIds(new Criteria(), $c)
            ->firstId();

        $productRepository = $this->container->get('product.repository');

        $productRepository->create(
            [
                [
                    'productNumber' => 'FREE-PRODUCT',
                    'stock' => 1,
                    'name' => 'FREE Product',
                    'active' => true,
                    'price' => [
                        [
                            'currencyId' => Defaults::CURRENCY,
                            'gross' => 0,
                            'net' => 0,
                            'linked' => true,
                        ],
                    ],
                    'manufacturer' => ['name' => 'test'],
                    'taxId' => $taxId,
                    'visibilities' => [
                        [
                            'salesChannelId' => $salesChannelId,
                            'visibility' => ProductVisibilityDefinition::VISIBILITY_ALL
                        ],
                    ],
                    'categories' => [
                        ['id' => $categoryId],
                    ],
                ],
            ],
            Context::createDefaultContext()
        );
    }
}
