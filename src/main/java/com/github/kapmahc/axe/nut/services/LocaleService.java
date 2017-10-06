package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.models.Locale;
import com.github.kapmahc.axe.nut.repositories.LocaleRepository;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import javax.annotation.Resource;

@Service("nut.localeService")
@Transactional(readOnly = true)
public class LocaleService {
    @CacheEvict(cacheNames = "locales", key = "#locale.#code")
    public void set(java.util.Locale locale, String code, String message) {
        String lang = locale2lang(locale);
        Locale it = localeRepository.findByLangAndCode(lang, code);
        if (it == null) {
            it = new Locale();
            it.setCode(code);
            it.setLang(lang);
        }
        it.setMessage(message);
        localeRepository.save(it);
    }

    @Cacheable(cacheNames = "locales", key = "#locale.#code")
    public String get(java.util.Locale locale, String code) {
        Locale it = localeRepository.findByLangAndCode(locale2lang(locale), code);
        return it == null ? null : it.getMessage();
    }

    private String locale2lang(java.util.Locale locale) {
        return locale.toLanguageTag();
    }

    @Resource
    LocaleRepository localeRepository;
}
