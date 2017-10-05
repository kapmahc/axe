package com.github.kapmahc.axe.ops.mail.repositories;

import com.github.kapmahc.axe.ops.mail.models.Domain;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.mail.domainRepository")
public interface DomainRepository extends CrudRepository<Domain, Long> {
}
