package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Locale;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.localeRepository")
public interface LocaleRepository extends CrudRepository<Locale, Long> {
    Locale findByLangAndCode(String lang, String code);
}
