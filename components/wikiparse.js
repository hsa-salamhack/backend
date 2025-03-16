import wtf from "wtf_wikipedia";
import wtfMarkdown from "wtf-plugin-markdown";
wtf.extend(wtfMarkdown);
import bionicifyMarkdown from "bionic-markdown";

async function parseBody(term, lang = "en") {
  let wikia = await wtf.fetch(term, lang);
  wikia = wikia
    .markdown()
    .replaceAll(/\[([^\]]+)\]\(\.\/(.*?)\)/g, `[$1](wiki/${lang}/$2)`)
    .replace(/(\|\s*.*\|\s*\n){2,}/g, "")
    .replace(
      /(\|.*\|[ \t]*\n)(\|[ \t]*---[ \t]*\|.*\n)?(\|.*\|[ \t]*\n)+/g,
      "",
    );

  return wikia;
}

async function getInfobox(term, lang = "en") {
  let wikia = await wtf.fetch(term, lang);
  let box = wikia.infobox().json();

  for (let key of Object.keys(box)) {
    if (/\.(png|jpe?g|gif|svg)$/i.test(box[key].text)) {
      let res = await fetch(
        `https://en.wikipedia.org/w/api.php?action=query&format=json&prop=imageinfo&iiprop=url&titles=File:${box[key].text}&origin=*`,
      );
      let data = await res.json();
      box[key] =
        Object.values(data.query.pages)[0]?.imageinfo?.[0]?.url || null;
    }
  }

  return box;
}

async function sections(term, lang = "en") {
  let wikia = await wtf.fetch(term, lang);
  let sections = {};
  let map = wikia.sections();

  map.forEach((section) => {
    if (section.title() && section.title() !== "") {
      sections[section.title()] = {
        body: section
          .markdown()
          .replaceAll(/\[([^\]]+)\]\(\.\/(.*?)\)/g, `[$1](wiki/${lang}/$2)`),
        body_bionic: bionicifyMarkdown(
          section
            .markdown()
            .replaceAll(/\[([^\]]+)\]\(\.\/(.*?)\)/g, `[$1](wiki/${lang}/$2)`),
        ),
      };
    }
  });

  return sections;
}

async function wiki(term, lang = "en") {
  let fullbody = await parseBody(term, lang);
  let bionic = await bionicifyMarkdown(fullbody);

  let summary = (fullbody.match(/^([\s\S]*?)(?=#+ |$)/)?.[1] || "").trim();
  let bionic_sum = await bionicifyMarkdown(summary);

  let sects = await sections(term, lang);
  let box = await getInfobox(term, lang);

  return {
    summary,
    infobox: box,
    bionic_summary: bionic_sum,
    sections: sects,

    // full_body_markdown: fullbody,
    // full_body_bionic: bionic,
  };
}

export { parseBody, sections, getInfobox as infobox, wiki };
