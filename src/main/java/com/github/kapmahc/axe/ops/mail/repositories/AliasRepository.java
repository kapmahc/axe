package com.github.kapmahc.axe.ops.mail.repositories;

import com.github.kapmahc.axe.ops.mail.models.Alias;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.mail.aliasRepository")
public interface AliasRepository extends CrudRepository<Alias, Long> {
}
