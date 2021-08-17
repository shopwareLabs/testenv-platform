import * as semver from "https://deno.land/x/semver/mod.ts";

async function main() {
    let releases = await getReleases();

    const ghConfig = {
        'fail-fast': false,
        matrix: {
            include: [] as any
        }
    };

    // Build
    for (let release of releases) {
        for (let tag of release.tags) {
            for (let imageName of release.imageNames) {
                ghConfig.matrix.include.push({
                    name: `Shopware ${tag}`,
                    runs: {
                        build: `cd images/6; docker buildx build --platform linux/amd64 --build-arg SHOPWARE_DL=${release.download} --build-arg SHOPWARE_VERSION=${release.version} --tag ${imageName}:${tag} --push .`
                    }
                });
            }
        }
    }

    await Deno.stdout.write(new TextEncoder().encode(JSON.stringify(ghConfig)));
}

function getMajorVersion(version: string) {
    let majorVersion = /\d+\.\d+.\d+/gm.exec(version);

    if (majorVersion && majorVersion[0]) {
        return majorVersion[0];
    } 

    return '';
}

main();

async function getReleases() {
    let json = await (await fetch('https://update-api.shopware.com/v1/releases/install?major=6&channel=rc')).json();
    let releases = [];
    let givenTags: string[] = [];


    for (let release of json) {
        const majorVersion = getMajorVersion(release.version);

        try {
            if (semver.lt(majorVersion, '6.4.0')) {
                continue;
            }
        } catch (e) {
        }

        if (!givenTags.includes(majorVersion)) {
            release.version = majorVersion;
            givenTags.push(majorVersion);
        } else {
            continue;
        }

        let image = {
            imageNames: ['shopware/testenv', 'ghcr.io/shopwareLabs/testenv'],
            version: release.version,
            download: release.uri,
            tags: [release.version]
        }

        releases.push(image);
    }

    return releases;
}
