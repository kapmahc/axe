package com.github.kapmahc.axe.forum.repositories;

import com.github.kapmahc.axe.forum.models.Comment;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("forum.commentRepository")
public interface CommentRepository extends CrudRepository<Comment, Long> {
}
