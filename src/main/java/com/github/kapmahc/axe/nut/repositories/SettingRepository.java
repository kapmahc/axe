package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Setting;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.settingRepository")
public interface SettingRepository extends CrudRepository<Setting, Long> {
}
