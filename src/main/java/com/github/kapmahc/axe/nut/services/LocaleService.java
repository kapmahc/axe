package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.models.Locale;
import com.github.kapmahc.axe.nut.repositories.LocaleRepository;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;

@Service("nut.localeService")
public class LocaleService {
    public void set(String lang, String code, String message) {

    }

    public String get(String lang, String code) {
        Locale it = localeRepository.findByLangAndCode(code, lang);
        return it == null ? null : it.getMessage();
    }

    @Resource
    LocaleRepository localeRepository;
}
