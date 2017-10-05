package com.github.kapmahc.axe.ops.mail.repositories;


import com.github.kapmahc.axe.ops.mail.models.User;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.mail.Repository")
public interface UserRepository extends CrudRepository<User, Long> {
}
