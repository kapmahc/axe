package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.models.Locale;
import com.github.kapmahc.axe.nut.repositories.LocaleRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import javax.annotation.Resource;

@Service("nut.localeService")
@Transactional(readOnly = true)
public class LocaleService {
    public void set(java.util.Locale locale, String code, String message) {
        String lang = locale2lang(locale);
        Locale it = localeRepository.findByLangAndCode(code, lang);
        if (it == null) {
            it = new Locale();
            it.setCode(code);
            it.setLang(lang);
        }
        it.setMessage(message);
        localeRepository.save(it);
    }

    public String get(java.util.Locale locale, String code) {
        Locale it = localeRepository.findByLangAndCode(code, locale2lang(locale));
        return it == null ? null : it.getMessage();
    }

    private String locale2lang(java.util.Locale locale) {
        return locale.toLanguageTag();
    }

    @Resource
    LocaleRepository localeRepository;
}
