package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Policy;
import com.github.kapmahc.axe.nut.models.Role;
import com.github.kapmahc.axe.nut.models.User;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository("nut.policyRepository")
public interface PolicyRepository extends CrudRepository<Policy, Long> {
    Policy findByUserAndRole(User user, Role role);

    List<Policy> findByUser(User user);
}
