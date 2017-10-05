package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.User;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.userRepository")
public interface UserRepository extends CrudRepository<User,Long> {
    User findByProviderTypeAndProviderId(User.Type type, String pid);
    User findByUid(String uid);
}
