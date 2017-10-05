package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Policy;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.policyRepository")
public interface PolicyRepository extends CrudRepository<Policy, Long> {
}
