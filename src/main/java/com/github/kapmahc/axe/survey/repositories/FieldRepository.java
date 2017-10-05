package com.github.kapmahc.axe.survey.repositories;

import com.github.kapmahc.axe.survey.models.Field;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.survey.FieldRepository")
public interface FieldRepository extends CrudRepository<Field, Long> {
}
