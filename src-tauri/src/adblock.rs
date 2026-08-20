use adblock::{Engine, FilterSet};
use adblock::lists::ParseOptions;

thread_local! {
    static BLOCKER: Engine = {
        let mut fs = FilterSet::new(false);
        let easylist = include_str!("../easylist.txt");
        let easyprivacy = include_str!("../easyprivacy.txt");
        fs.add_filters(easylist.lines(), ParseOptions::default());
        fs.add_filters(easyprivacy.lines(), ParseOptions::default());
        Engine::from_filter_set(fs, true)
    };
}

pub fn get_cosmetic_css(url: &str) -> String {
    BLOCKER.with(|blocker| {
        let resources = blocker.url_cosmetic_resources(url);
        let mut css = String::new();
        for hide in &resources.hide_selectors {
            css.push_str(hide);
            css.push_str("{display:none!important}\n");
        }
        css
    })
}

pub fn get_injected_script(url: &str) -> String {
    BLOCKER.with(|blocker| {
        blocker.url_cosmetic_resources(url).injected_script.clone()
    })
}
