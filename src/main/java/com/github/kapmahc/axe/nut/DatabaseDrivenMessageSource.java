package com.github.kapmahc.axe.nut;

import com.github.kapmahc.axe.nut.services.LocaleService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ResourceLoaderAware;
import org.springframework.context.support.AbstractMessageSource;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.core.io.ResourceLoader;
import org.springframework.lang.Nullable;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import javax.annotation.Resource;
import java.text.MessageFormat;
import java.util.Locale;

@Component("messageSource")
public class DatabaseDrivenMessageSource extends AbstractMessageSource implements ResourceLoaderAware {
    @Override
    public void setResourceLoader(ResourceLoader resourceLoader) {

    }

    @Nullable
    @Override
    protected MessageFormat resolveCode(String code, Locale locale) {
        String msg = tr(code, locale);
        return createMessageFormat(msg, locale);
    }

    @Nullable
    @Override
    protected String resolveCodeWithoutArguments(String code, Locale locale) {
        return tr(code, locale);
    }

    @PostConstruct
    void init() {
        ReloadableResourceBundleMessageSource rbm = new ReloadableResourceBundleMessageSource();
        rbm.setBasename(basename);
        rbm.setDefaultEncoding(encoding);
        rbm.setCacheSeconds(cacheSeconds);
        rbm.setFallbackToSystemLocale(fallbackToSystemLocale);
        setParentMessageSource(rbm);
    }

    private String tr(String code, Locale locale) {
        return localeService.get(locale, code);
    }

    @Value("${spring.messages.basename}")
    String basename;
    @Value("${spring.messages.encoding}")
    String encoding;
    @Value("${spring.messages.cache-seconds}")
    int cacheSeconds;
    @Value("${spring.messages.fallback-to-system-locale}")
    boolean fallbackToSystemLocale;

    @Resource
    LocaleService localeService;

    private final static Logger logger = LoggerFactory.getLogger(DatabaseDrivenMessageSource.class);
}
