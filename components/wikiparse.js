import wtf from "wtf_wikipedia";
import wtfMarkdown from "wtf-plugin-markdown";
wtf.extend(wtfMarkdown);
import bionicifyMarkdown from "bionic-markdown";

async function parseBody(wiki, lang = "en") {
  let wikia = wiki
    .markdown()
    .replaceAll(/\[([^\]]+)\]\(\.\/(.*?)\)/g, `[$1](wiki/${lang}/$2)`)
    .replace(/(\|\s*.*\|\s*\n){2,}/g, "")
    .replace(
      /(\|.*\|[ \t]*\n)(\|[ \t]*---[ \t]*\|.*\n)?(\|.*\|[ \t]*\n)+/g,
      "",
    );

  return wikia;
}

async function getImages(wikia, lang = "en") {
  if (!wikia.infobox())
    return {
      firstImage: wikia.images()[0].url(),
    };
  let box = wikia.infobox();

  let boxj = box.json();
  let imgs = {};

  for (let key of Object.keys(boxj)) {
    if (/\.(png|jpe?g|gif|svg)$/i.test(boxj[key].text)) {
      let res = await fetch(
        `https://${lang}.wikipedia.org/w/api.php?action=query&format=json&prop=imageinfo&iiprop=url&titles=File:${boxj[key].text}&origin=*`,
      );
      let data = await res.json();
      imgs[key] =
        Object.values(data.query.pages)[0]?.imageinfo?.[0]?.url || null;
    }
  }
  imgs.firstImage = wikia.images()[0].url();
  return imgs;
}

async function sections(wikia, lang) {
  let sections = {};
  let map = wikia.sections();

  const sectionPromises = Array.from(map).map(async (section) => {
    if (section.title() && section.title() !== "") {
      const bodyMarkdown = section
        .markdown()
        .replaceAll(/\[([^\]]+)\]\(\.\/(.*?)\)/g, `[$1](wiki/${lang}/$2)`);

      const bodyBionic = await bionicifyMarkdown(bodyMarkdown);

      return {
        title: section.title(),
        data: {
          body: bodyMarkdown,
          body_bionic: bodyBionic,
        },
      };
    }
    return null;
  });

  let results = await Promise.all(sectionPromises);

  results.forEach((result) => {
    if (result) {
      sections[result.title] = result.data;
    }
  });

  return sections;
}

async function wiki(term, lang = "en") {
  let wikia = await wtf.fetch(term, lang);

  let [fullbody, sects, imgs] = await Promise.all([
    parseBody(wikia, lang),
    sections(wikia, lang),
    getImages(wikia, lang),
  ]);

  let [bionic, summary] = await Promise.all([
    bionicifyMarkdown(fullbody),
    (fullbody.match(/^([\s\S]*?)(?=#+ |$)/)?.[1] || "").trim(),
  ]);

  let bionic_sum = await bionicifyMarkdown(summary);

  return {
    summary,
    images: imgs,
    thumbnail: imgs?.firstImage || null,
    bionic_summary: bionic_sum,
    sections: sects,
    full_body: fullbody,
  };
}

export { wiki };
