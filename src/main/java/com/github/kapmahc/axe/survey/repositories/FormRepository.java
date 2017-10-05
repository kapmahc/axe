package com.github.kapmahc.axe.survey.repositories;

import com.github.kapmahc.axe.survey.models.Form;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.survey.FormRepository")
public interface FormRepository  extends CrudRepository<Form, Long> {
}
