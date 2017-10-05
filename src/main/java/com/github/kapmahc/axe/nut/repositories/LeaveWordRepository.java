package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.LeaveWord;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.leaveWordRepository")
public interface LeaveWordRepository extends CrudRepository<LeaveWord, Long> {
}
