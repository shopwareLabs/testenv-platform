import * as semver from "https://deno.land/x/semver/mod.ts";

async function main() {
    let releases = await getReleases();

    console.log(releases);

    // Build
    for (let release of releases) {
        for (let tag of release.tags) {
            for (let imageName of release.imageNames) {
                let process = Deno.run({
                    cmd: ['docker', 'build', '-t', `${imageName}:${tag}`, '--build-arg', `SHOPWARE_DL=${release.download}`, '.'],
                    stdout: 'inherit'
                });
    
                const {success} = await process.status();
    
                if (!success) {
                    Deno.exit(-1);
                }
            }
        }
    }

    // Push
    for (let release of releases) {
        for (let tag of release.tags) {
            for (let imageName of release.imageNames) {
                let process = Deno.run({
                    cmd: ['docker', 'push', `${imageName}:${tag}`],
                    stdout: 'inherit'
                });

                const {success} = await process.status();

                if (!success) {
                    Deno.exit(-1);
                }
            }
        }
    }
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
    let json = await (await fetch('https://update-api.shopware.com/v1/releases/install?major=6')).json();
    let releases = [];
    let givenTags: string[] = [];


    for (let release of json) {
        try {
            if (semver.lt(release.version, '6.3.0')) {
                continue;
            }
        } catch (e) {
        }

        const majorVersion = getMajorVersion(release.version);

        if (!givenTags.includes(majorVersion)) {
            release.version = majorVersion;
            givenTags.push(majorVersion);
        } else {
            continue;
        }

        let image = {
            imageNames: ['shopware/testenv'],
            version: release.version,
            download: release.uri,
            tags: [release.version]
        }

        releases.push(image);
    }

    return releases;
}
