package com.github.kapmahc.nut.dao;

import org.skife.jdbi.v2.sqlobject.Bind;
import org.skife.jdbi.v2.sqlobject.SqlQuery;
import org.skife.jdbi.v2.sqlobject.SqlUpdate;

import java.util.Date;

public interface LocaleDao {
    @SqlUpdate("insert into locales (lang, code, message, updatedAt) values (:lang, :code, :message, :updatedAt)")
    void insert(@Bind("lang") String lang, @Bind("code") String code, @Bind("message") String message, @Bind("updatedAt") Date updatedAt);

    @SqlUpdate("update locales set message=:message, updatedAt=:updatedAt where lang = :lang and code = :code")
    void update(@Bind("lang") String lang, @Bind("code") String code, @Bind("message") String message, @Bind("updatedAt") Date updatedAt);

    @SqlQuery("select message from locales where lang = :lang and code = :code")
    String select(@Bind("lang") String lang, @Bind("code") String code);
}
