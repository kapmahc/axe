package com.github.kapmahc.axe.survey.repositories;

import com.github.kapmahc.axe.survey.models.Record;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.survey.RecordRepository")
public interface RecordRepository  extends CrudRepository<Record, Long> {
}
