package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Log;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.logRepository")
public interface LogRepository extends CrudRepository<Log, Long> {
}
