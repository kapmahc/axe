package com.github.kapmahc.axe.forum.repositories;

import com.github.kapmahc.axe.forum.models.Article;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("forum.articleRepository")
public interface ArticleRepository extends CrudRepository<Article, Long> {
}
